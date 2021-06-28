import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { map } from "rxjs/operators";
import { shareReplay } from "rxjs/operators";
import { ApiPaths } from "src/environments/paths";
import { Observable, of, throwError } from "rxjs";
import { environment } from "../../environments/environment";

class Token {
  accessToken: string = "";
  refreshToken: string = "";
  expirationAccess: Date = new Date();
  expirationRefresh: Date = new Date();
}

// interface httpHeaders {
//   headers: {
//     "Content-Type": string;
//     "Authorization": string;
//   };
// }

@Injectable({
  providedIn: "root",
})
export class AuthService {
  private url = `${environment.url}/${ApiPaths.auth}`;
  private urlRefresh = `${environment.url}/${ApiPaths.refreshToken}`;
  private urlLogout = `${environment.url}/${ApiPaths.logout}`;
  private opts = {
    headers: {
      "Content-Type": "application/json",
    },
  };
  constructor(private http: HttpClient) {}

  auth(username: string, password: string): Observable<Token> {
    // async auth(username: string, password: string): Promise<any> {
    return this.http.post<Token>(this.url, { username, password }, this.opts)
      .pipe(
        map(token => {
          localStorage.setItem("token", JSON.stringify(token));
          return token;
        }),
      );
  }

  getTokenFromLS(): Token {
    let tkn = new Token();
    const tknStr = localStorage.getItem("token");
    if (!tknStr) {
      console.error("No hay token en el localStorage");
    } else if (typeof tknStr === "string") {
      tkn = JSON.parse(tknStr);
    }
    return tkn;
  }

  refreshToken(tkn: Token): Observable<Token> {
    const headers = {
      headers: {
        "Content-Type": "application/json",
        "Refresh": `${tkn.refreshToken}`,
      },
    };
    return this.http.get<Token>(this.urlRefresh, headers)
      .pipe(
        map(token => {
          localStorage.setItem("token", JSON.stringify(token));
          return token;
        }),
        shareReplay(),
      );
  }

  validateToken(): Observable<Token> {
    const tkn = this.getTokenFromLS();
    if (!this.isLoggedIn()) {
      this.logout();
      // en lugar de lanzar error, redireccionar a login o home
      throw throwError("El Refresh Token es muy antiguo");
    } else if (this.isAccessTokenExpired()) {
      return this.refreshToken(tkn);
    } else {
      return of(tkn);
    }
  }

  isLoggedIn(): boolean {
    const expRefresh: Date = this.getExpirationRefresh();
    const now = new Date();
    return expRefresh.getTime() - now.getTime() > 5000;
  }

  isAccessTokenExpired(): boolean {
    const expAccess = this.getExpirationAccess();
    const now = new Date();
    return expAccess.getTime() - now.getTime() < 5000;
  }

  getExpirationAccess(): Date {
    const token = this.getTokenFromLS();
    return new Date(token.expirationAccess);
  }

  getExpirationRefresh(): Date {
    const token = this.getTokenFromLS();
    return new Date(token.expirationRefresh);
  }

  logout(): void {
    this.http.get(this.urlLogout).subscribe(
      _success => localStorage.removeItem("token"),
      err => console.error(err),
    );
  }
}
