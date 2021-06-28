import { Component, OnInit } from "@angular/core";
import { mergeMap } from "rxjs/operators";
import { Usuario } from "src/app/models/usuario";
import { AuthService } from "src/app/services/auth.service";
import { UsuariosService } from "src/app/services/usuarios.service";

@Component({
  selector: "app-usuarios",
  templateUrl: "./usuarios.component.html",
  styleUrls: ["./usuarios.component.css"],
})
export class UsuariosComponent implements OnInit {
  constructor(private authService: AuthService, private usuarioService: UsuariosService) {}

  ngOnInit(): void {
    const usuario = new Usuario(1, "prueba", "prueba", "nada@ts.s", true);
    this.authService.auth(usuario.username, usuario.password)
      .pipe(
        mergeMap((_token, _n) => {
          return this.usuarioService.getUsuarios();
        }),
      ).subscribe(data => console.log(data));
  }
}
