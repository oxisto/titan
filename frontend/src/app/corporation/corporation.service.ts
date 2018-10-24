import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import 'rxjs/add/operator/map';
import { Observable } from 'rxjs/Observable';
import { AuthService } from '../auth/auth.service';
import { Corporation } from './corporation';



@Injectable()
export class CorporationService {

  corporation: Observable<Corporation>;

  constructor(private http: HttpClient,
    private authService: AuthService) {
    if (authService.isLoggedIn()) {
      this.fetch();
    }
  }

  get(): Observable<Corporation> {
    return this.http.get<Corporation>('/api/corporation');
  }

  fetch() {
    this.corporation = this.get();
  }

  getCorporationLogo(corporationID: number) {
    return 'https://image.eveonline.com/Corporation/' + corporationID + '_64.png';
  }

}
