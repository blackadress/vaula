import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { mergeMap } from "rxjs/operators";
import { environment } from "src/environments/environment";
import { ApiPaths } from "src/environments/paths";
import { Observable } from "rxjs";
import { Usuario } from "../models/usuario";
import { AuthService } from "./auth.service";

@Injectable({
  providedIn: "root",
})
export class UsuariosService {
  private usuarios: Usuario[];
  private url = `${environment.url}/${ApiPaths.users}`;
  private headers = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  constructor(private http: HttpClient, private authService: AuthService) {
    this.usuarios = [];
  }

  createUsuario(usuario: Usuario): Observable<Usuario> {
    return this.http.post<Usuario>(this.url, usuario, this.headers);
  }

  getUsuarios(): Observable<Usuario[]> {
    return this.authService.validateToken()
      .pipe(
        mergeMap((token, n) => {
          const header = {
            headers: {
              "Content-Type": "application/json",
              "Authorization": `Bearer ${token.accessToken}`,
            },
          };
          return this.http.get<Usuario[]>(this.url, header);
        }),
      );
  }
}
