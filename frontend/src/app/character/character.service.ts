import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Character } from './character';

@Injectable()
export class CharacterService {

  constructor(private http: HttpClient) {
  }

  getCharacter(): Observable<Character> {
    return this.http.get<Character>('/api/character');
  }

  getCharacterPortraitURL(characterID: number, size: number) {
    return 'https://image.eveonline.com/Character/' + characterID + '_' + size + '.jpg';
  }

}
