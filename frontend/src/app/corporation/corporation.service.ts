import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Corporation, Wallets } from './corporation';

@Injectable()
export class CorporationService {

  constructor(private http: HttpClient) {
  }

  getCorporation(): Observable<Corporation> {
    return this.http.get<Corporation>('/api/corporation');
  }

  getCorporationWallets(): Observable<Wallets> {
    return this.http.get<Wallets>('/api/corporation/wallets');
  }

  getCorporationLogo(corporationID: number) {
    return 'https://image.eveonline.com/Corporation/' + corporationID + '_64.png';
  }

}
