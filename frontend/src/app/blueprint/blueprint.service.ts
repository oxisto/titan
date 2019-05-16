import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { AuthService } from '../auth/auth.service';
import { Blueprint } from './blueprint';

@Injectable()
export class BlueprintService {

  constructor(private http: HttpClient,
    private authService: AuthService) {

  }

  getBlueprints(): Observable<Blueprint[]> {
    return this.http.get<Blueprint[]>('/api/blueprints');
  }

  getBlueprint(typeID: number) {
    return this.http.get<Blueprint>('/api/blueprints/' + typeID);
  }

  getManufacturing(typeID: number, ME: number, TE: number, facilityTax: number) {
    return this.http.get('/api/manufacturing/' + typeID + '?ME=' + ME + '&TE=' + TE + '&facilityTax=' + facilityTax);
  }

}
