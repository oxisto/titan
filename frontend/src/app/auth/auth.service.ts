import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { HttpParams } from '@angular/common/http';
import { JwtHelperService } from '@auth0/angular-jwt';

const helper = new JwtHelperService();

@Injectable()
export class AuthService {

  constructor(private router: Router) {
    const params = new HttpParams({ fromString: window.location.hash.replace('#?', '') });

    const idToken = params.get('token');

    if (idToken) {
      console.log('Setting access token in localStorage...');
      localStorage.setItem('token', idToken);
      this.router.navigate(['/']);
    }
  }

  isLoggedIn() {
    const token = localStorage.getItem('token');

    return !helper.isTokenExpired(token);
  }

}
