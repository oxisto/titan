import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { AuthService } from '../auth/auth.service';
import { Character } from './character';

@Injectable()
export class CharacterService {

  character: Observable<Character>;

  constructor(private http: HttpClient,
    private authService: AuthService) {
    this.fetch();
  }

  get(): Observable<Character> {
    return this.http.get<Character>('/api/character');
  }

  fetch() {
    this.character = this.get();
  }

  getCharacterPortraitURL(characterID: number, size: number) {
    return 'https://image.eveonline.com/Character/' + characterID + '_' + size + '.jpg';
  }

}
