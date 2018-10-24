import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpClient } from '@angular/common/http';

import 'rxjs/add/operator/map';

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

  getManufacturing(typeID: number) {
    return this.http.get<any>('/api/manufacturing/' + typeID);
  }

}
