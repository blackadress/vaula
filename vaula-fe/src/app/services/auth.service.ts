import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { map } from "rxjs/operators";
import { ApiPaths } from "src/environments/paths";
import { Observable, of, throwError } from "rxjs";

class Token {
  accessToken: string = "";
  refreshToken: string = "";
  expirationAccess: Date = new Date();
  expirationRefresh: Date = new Date();
}

@Injectable({
  providedIn: "root",
})
export class AuthService {
  private url = "http://localhost:8000/api/token";
  private urlRefresh = "http://localhost:8000/api/refresh";
  private opts = {
    headers: {
      "Content-Type": "application/json",
    },
  };
  constructor(private http: HttpClient) {}

  auth(username: string, password: string): void {
    const authUser = {
      username: username,
      password: password,
    };
    this.http.post<Token>(this.url, authUser, this.opts)
      .subscribe(token => {
        localStorage.setItem("token", JSON.stringify(token));
      });
  }

  getTokenFromLS(): Token {
    let tkn = new Token();
    const tknStr = localStorage.getItem("token");
    if (tknStr) {
      console.error("No hay token en el localStorage");
    } else if (typeof tknStr === "string") {
      tkn = JSON.parse(tknStr);
    }
    return tkn;
  }

  refreshToken(tkn: Token): Observable<Token> {
    return this.http.post<Token>(this.urlRefresh, tkn, this.opts)
      .pipe(
        map(token => {
          localStorage.setItem("token", JSON.stringify(token));
          return token;
        }),
      );
  }

  refreshAccessToken(): Observable<Token> {
    const tkn = this.getTokenFromLS();
    const now = new Date();
    const segundos_5 = 5000;
    if (tkn.expirationRefresh.getTime() - now.getTime() > segundos_5) {
      this.logout();
      throw throwError("El Refresh Token es muy antiguo");
    } else if (tkn.expirationAccess.getTime() - now.getTime() < segundos_5) {
      return this.refreshToken(tkn);
    } else if (tkn.expirationAccess.getTime() - now.getTime() > segundos_5) {
      return of(tkn);
    } else {
      throw throwError("WTF");
    }
  }

  logout(): void {
    this.http.get(ApiPaths.logout).subscribe(
      _success => localStorage.clear(),
      err => console.error(err),
    );
  }
}
