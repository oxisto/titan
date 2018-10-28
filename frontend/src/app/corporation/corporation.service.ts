import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Corporation } from './corporation';

@Injectable()
export class CorporationService {

  constructor(private http: HttpClient) {
  }

  getCorporation(): Observable<Corporation> {
    return this.http.get<Corporation>('/api/corporation');
  }

  getCorporationLogo(corporationID: number) {
    return 'https://image.eveonline.com/Corporation/' + corporationID + '_64.png';
  }

}
