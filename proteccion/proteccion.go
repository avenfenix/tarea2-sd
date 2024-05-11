package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	pb "github.com/tarea2/proteccion/proto"
	"google.golang.org/grpc"
)

const (
	APIEntryPoint1 = "https://api.ilovepdf.com/v1/auth"
	APIEntryPoint2 = "https://api.ilovepdf.com/v1/start"
)

type server struct {
	pb.UnimplementedProtectorServer
}

type ILovePdf struct {
	PublicKey string
	Token     string
}

func NewILovePdf(publicKey string) *ILovePdf {
	resp, _ := http.PostForm(APIEntryPoint1, map[string][]string{
		"public_key": {publicKey},
	})
	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	return &ILovePdf{PublicKey: publicKey, Token: result["token"]}
}

type Operations struct {
	*ILovePdf
	Token  string
	TaskID string
	Tool   string
	Server string
	Files  []map[string]string
}

func NewOperations(publicKey string) *Operations {
	op := &Operations{ILovePdf: NewILovePdf(publicKey)}
	op.retrieveToken()
	return op
}

// Metodo para guardar el token
func (op *Operations) retrieveToken() {
	if op.Token == "" {
		resp, err := http.PostForm(APIEntryPoint1, map[string][]string{
			"public_key": {op.PublicKey},
		})
		if err != nil {
			return
		}
		var result map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return
		}
		resp.Body.Close()
		op.Token = result["token"]
	}
}

func (op *Operations) startTask(tool string) {
	op.Tool = tool
	op.retrieveToken() //Recibir token
	req, _ := http.NewRequest("GET", APIEntryPoint2+"/"+tool, nil)
	req.Header.Set("Authorization", "Bearer "+op.Token)
	resp, _ := http.DefaultClient.Do(req)
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	op.TaskID, op.Server = result["task"].(string), result["server"].(string)
}

func (op *Operations) addFile(filename string) error {
	// Verificar si el archivo existe
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("el archivo especificado no existe: %s", filename)
	}

	// Construir la URL de la solicitud
	url := fmt.Sprintf("https://%s/v1/upload", op.Server)

	// Abrir el archivo
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	// Crear un buffer para el cuerpo del formulario
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Agregar el archivo al formulario
	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return fmt.Errorf("error al crear la parte del formulario: %v", err)
	}
	if _, err = io.Copy(part, file); err != nil {
		return fmt.Errorf("error al copiar el contenido del archivo: %v", err)
	}

	// Agregar el parámetro "task" al formulario
	writer.WriteField("task", op.TaskID)

	// Cerrar el escritor multipart
	if err := writer.Close(); err != nil {
		return fmt.Errorf("error al cerrar el escritor multipart: %v", err)
	}

	// Crear la solicitud HTTP POST
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("error al crear la solicitud http: %v", err)
	}

	// Establecer el tipo de contenido en la solicitud
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Agregar el encabezado de autorización
	req.Header.Set("Authorization", "Bearer "+op.Token)

	// Realizar la solicitud HTTP
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error al realizar la solicitud http: %v", err)
	}
	defer resp.Body.Close()

	// Decodificar la respuesta JSON
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("error al decodificar la respuesta json: %v", err)
	}

	// Verificar si el archivo se agregó correctamente
	if serverFilename, ok := response["server_filename"].(string); ok {
		op.Files = append(op.Files, map[string]string{
			"server_filename": serverFilename,
			"filename":        filename,
		})
		return nil
	}

	return fmt.Errorf("error al agregar el archivo: %v", response)
}

func (op *Operations) execute(password string) {
	url := fmt.Sprintf("https://%s/v1/process", op.Server)
	params := map[string]interface{}{
		"task":     op.TaskID,
		"tool":     op.Tool,
		"files":    op.Files,
		"password": password,
	}
	jsonData, _ := json.Marshal(params)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+op.Token)
	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud HTTP
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error executing task:", err)
		return
	}
	defer resp.Body.Close()

}

func (op *Operations) download(outputFilename string, inputPath string) string {
	url := fmt.Sprintf("https://%s/v1/download/%s", op.Server, op.TaskID)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+op.Token)
	resp, _ := http.DefaultClient.Do(req)

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.Status) // Imprime cualquier mensaje de error en la respuesta
		return ""
	}

	// Obtener el directorio del archivo de entrada
	outputDir := filepath.Dir(inputPath)

	// Concatenar el directorio y el nombre de archivo de salida
	outputPath := filepath.Join(outputDir, outputFilename)

	out, _ := os.Create(outputPath)
	defer out.Close()
	io.Copy(out, resp.Body)
	resp.Body.Close()

	return outputPath
}

func (s *server) Proteger(c context.Context, in *pb.SolicitudProteger) (*pb.RespuestaProteger, error) {
	log.Printf("Peticion recibida con correo: %v", in.GetCorreo())

	// Guardar archivo enviado por el cliente en carpeta ./files

	archivo := in.GetFile()

	// Guardar el archivo temporalmente en el servidor
	archivoTemp, err := ioutil.TempFile("", "archivo_")
	if err != nil {
		return nil, fmt.Errorf("error al crear archivo temporal: %v", err)
	}
	defer archivoTemp.Close()

	// Escribir el contenido del archivo en el archivo temporal
	_, err = archivoTemp.Write(archivo)
	if err != nil {
		return nil, fmt.Errorf("error al escribir en archivo temporal: %v", err)
	}

	// Obtener la ruta del archivo temporal
	path := archivoTemp.Name()

	// CHATGPT aqui guarda el archivo

	password := in.GetRut()

	publicKey := os.Getenv("PUBLIC_KEY")
	op := NewOperations(publicKey)
	op.startTask("protect")
	op.addFile(path)
	op.execute(password)
	fileName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) + "_protegido.pdf"
	targetPath := op.download(fileName, path)
	if targetPath == "" {
		return &pb.RespuestaProteger{Message: "Error al proteger el archivo"}, nil
	}

	return &pb.RespuestaProteger{Message: "Solicitud de proteccion recibida y sera enviada a: " + in.GetCorreo()}, nil
}

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error al leer el archivo .env")
	}

	listener, err := net.Listen("tcp", ":"+os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterProtectorServer(s, &server{})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
