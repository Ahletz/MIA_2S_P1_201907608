package Tools

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Estructura para recibir los datos desde el frontend
type CommandRequest struct {
	Command string `json:"command"`
}

// Estructura para enviar la respuesta al frontend
type CommandResponse struct {
	Result string `json:"resultado"`
}

// Handler que recibe el texto desde el frontend, lo procesa y envía la respuesta
func handleCommand(w http.ResponseWriter, r *http.Request) {
	// Asegúrate de permitir solicitudes de dominios diferentes (CORS)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// Leer el cuerpo de la solicitud
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "No se pudo leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// Deserializar el JSON en la estructura CommandRequest
	var commandRequest CommandRequest
	err = json.Unmarshal(body, &commandRequest)
	if err != nil {
		http.Error(w, "Error al deserializar el JSON", http.StatusBadRequest)
		return
	}

	// Procesar el comando usando tu función
	result := Separar_cmd(commandRequest.Command)

	// Crear la respuesta
	commandResponse := CommandResponse{
		Result: result,
	}

	// Serializar la respuesta como JSON y enviarla de vuelta
	response, err := json.Marshal(commandResponse)
	if err != nil {
		http.Error(w, "Error al serializar la respuesta", http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

// Función para iniciar el servidor
func StartServer() {
	http.HandleFunc("/api/command", handleCommand)
	log.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
