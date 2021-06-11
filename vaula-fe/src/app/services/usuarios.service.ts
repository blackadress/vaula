import { Injectable } from "@angular/core";
import { Usuario } from "../models/usuario";

@Injectable({
  providedIn: "root",
})
export class UsuariosService {
  private usuarios: Usuario[];

  constructor() {
    this.usuarios = [];
  }

  getUsuarios(): Usuario[] {
    return this.usuarios;
  }
}
