import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { AuthService } from '../auth/auth.service';

@Injectable()
export class ManufacturingService {

  constructor(private http: HttpClient,
    private authService: AuthService) {

  }

  getManufacturingTypeIDs(options: {
    categoryIDs: number[],
    sortBy: string,
    nameFilter?: string,
    hasRequiredSkillsOnly?: boolean,
    maxProductionCosts?: number,
    metaGroupID?: number
  }) {
    let params = new HttpParams().set('sortBy', options.sortBy);

    if (options.metaGroupID) {
      params = params.set('metaGroupID', options.metaGroupID.toString());
    }

    if (options.nameFilter) {
      params = params.set('nameFilter', options.nameFilter);
    }

    if (options.maxProductionCosts) {
      params = params.set('maxProductionCosts', options.maxProductionCosts.toString());
    }

    params = params.set('hasRequiredSkillsOnly', String(options.hasRequiredSkillsOnly));
    params = params.set('categoryIDs', options.categoryIDs.join(','));

    return this.http.get<any[]>('/api/manufacturing', { params: params });
  }

  getManufacturing(typeID: number) {
    return this.http.get('/api/manufacturing/' + typeID);
  }

  getManufacturingCategories(): Observable<any[]> {
    return this.http.get<any[]>('/api/manufacturing-categories');
  }

}
