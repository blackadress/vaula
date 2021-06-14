import { Component, OnInit } from "@angular/core";
import { Usuario } from "src/app/models/usuario";
import { AuthService } from "src/app/services/auth.service";

@Component({
  selector: "app-usuarios",
  templateUrl: "./usuarios.component.html",
  styleUrls: ["./usuarios.component.css"],
})
export class UsuariosComponent implements OnInit {
  constructor(private authService: AuthService) {}

  ngOnInit(): void {
    const usuario = new Usuario(1, "user_test", "1234", "nada@ts.s", true);
    this.authService.auth(usuario.username, usuario.password);
    const tkn = this.authService.getTokenFromLS();
    console.log(tkn);
  }
}
