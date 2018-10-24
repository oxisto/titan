import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';



import 'rxjs/add/operator/map';

import { AuthService } from '../auth/auth.service';
import { Character } from './character';

@Injectable()
export class CharacterService {

  character: Observable<Character>;

  constructor(private http: HttpClient,
    private authService: AuthService) {
    if (authService.isLoggedIn()) {
      this.fetch();
    }
  }

  get(): Observable<Character> {
    return this.http.get<Character>('/api/character');
  }

  fetch() {
    this.character = this.get();
  }

  getCharacterPortraitURL(characterID: number) {
    return 'https://image.eveonline.com/Character/' + characterID + '_128.jpg';
  }

}
