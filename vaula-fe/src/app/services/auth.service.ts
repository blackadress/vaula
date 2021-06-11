import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { Usuario } from "../models/usuario";

class Token {
  accessToken: string = "";
  refreshToken: string = "";
}

@Injectable({
  providedIn: "root",
})
export class AuthService {
  private url = "http://localhost:8000/api/token";
  private opts = {
    headers: {
      "Content-Type": "application/json",
    },
  };
  constructor(private http: HttpClient) {}

  getToken(usuario: Usuario): Observable<Token> {
    const authUser = {
      username: usuario.username,
      password: usuario.password,
    };
    // return this.http.post<Token>(this.url, authUser, this.opts);
    return this.http.post<Token>(this.url, authUser, this.opts);
  }
}
