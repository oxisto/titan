import { Component, OnInit } from '@angular/core';
import { Character } from './character/character';
import { CharacterService } from './character/character.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  loggedIn = false;

  title = 'app works!';
  character: Character;
  characterPortrait: string;

  constructor(private characterService: CharacterService) {
    this.characterService.character.subscribe(character => {
      this.character = character;
      this.characterPortrait = this.characterService.getCharacterPortraitURL(character.characterID, 64);
    });
  }

  ngOnInit() {

  }

}
