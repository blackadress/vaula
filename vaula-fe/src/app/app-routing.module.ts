import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { AlumnosComponent } from "./components/alumnos/alumnos.component";
import { CursosComponent } from "./components/cursos/cursos.component";
import { ProfesoresComponent } from "./components/profesores/profesores.component";
import { UsuariosComponent } from "./components/usuarios/usuarios.component";

const routes: Routes = [
  { path: "usuarios", component: UsuariosComponent },
  { path: "profesores", component: ProfesoresComponent },
  { path: "alumnos", component: AlumnosComponent },
  { path: "cursos", component: CursosComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
