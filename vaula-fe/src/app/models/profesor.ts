import { Usuario } from "./usuario";

class Profesor {
  id: number;
  nombres: string;
  apellidos: string;
  usuarioId: number;
  usuario: Usuario;
  activo: boolean;
  createdAt: Date;
  updatedAt: Date;

  constructor(
    id: number,
    nombres: string,
    apellidos: string,
    usuarioId: number,
    usuario: Usuario,
    activo: boolean,
    createdAt: Date = new Date(),
    updatedAt: Date = new Date(),
  ) {
    this.id = id;
    this.nombres = nombres;
    this.apellidos = apellidos;
    this.usuarioId = usuarioId;
    this.usuario = usuario;
    this.activo = activo;
    this.createdAt = createdAt;
    this.updatedAt = updatedAt;
  }
}
