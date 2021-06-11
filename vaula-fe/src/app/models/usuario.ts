export class Usuario {
  id: number;
  username: string;
  password: string;
  email: string;
  activo: boolean;
  createdAt: Date;
  updatedAt: Date;

  constructor(
    id: number,
    username: string,
    password: string,
    email: string,
    activo: boolean,
    createdAt: Date = new Date(),
    updatedAt: Date = new Date(),
  ) {
    this.id = id;
    this.username = username;
    this.password = password;
    this.email = email;
    this.activo = activo;
    this.createdAt = createdAt;
    this.updatedAt = updatedAt;
  }
}
