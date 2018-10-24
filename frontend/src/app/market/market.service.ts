import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { AuthService } from '../auth/auth.service';

@Injectable()
export class MarketService {

  constructor(private http: HttpClient,
    private authService: AuthService) {

  }

  postOpenMarketView(typeID: number): Observable<any> {
    return this.http.post<any>('/api/market/view?typeID=' + typeID, null);
  }
}
