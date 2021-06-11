import { NgModule } from "@angular/core";
import { BrowserModule } from "@angular/platform-browser";

import { HttpClientModule } from "@angular/common/http";
import { AlumnosComponent } from "./components/alumnos/alumnos.component";
import { CursosComponent } from "./components/cursos/cursos.component";
import { ProfesoresComponent } from "./components/profesores/profesores.component";
import { UsuariosComponent } from "./components/usuarios/usuarios.component";
import { AppRoutingModule } from "./app-routing.module";
import { AppComponent } from "./app.component";

@NgModule({
  declarations: [
    AppComponent,
    AlumnosComponent,
    ProfesoresComponent,
    CursosComponent,
    UsuariosComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
  ],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
