import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';  // Asegúrate de importar HttpClient

@Component({
  selector: 'app-comando',
  templateUrl: './comando.component.html',
  styleUrls: ['./comando.component.css']  // Cambia "styleUrl" a "styleUrls"
})
export class ComandoComponent {
  inputText: string = '';  // Para almacenar el texto del cuadro de comandos
  responseText: string = '';  // Para almacenar la respuesta del backend

  constructor(private http: HttpClient) {}  // Inyecta HttpClient para hacer solicitudes HTTP

  // Este es el método en Angular para enviar el comando al backend en Go
sendText() {
  const apiUrl = 'http://localhost:8080/api/command';  // URL del servidor Go

  // Realizamos una petición POST al backend con el texto del input
  this.http.post<any>(apiUrl, { command: this.inputText }).subscribe(
    (response) => {
      this.responseText = response.resultado;  // Asignamos la respuesta al cuadro gris
    },
    (error) => {
      console.error('Error al enviar el comando:', error);
      this.responseText = 'Error al enviar el comando';  // Manejar error
    }
  );
}


  // Función para cargar un archivo .txt
  uploadFile(event: any) {
    const file = event.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (e: any) => {
        this.inputText = e.target.result;  // Carga el contenido del archivo al cuadro de texto
      };
      reader.readAsText(file);
    }
  }
}

