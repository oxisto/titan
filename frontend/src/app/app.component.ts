import { Component, OnInit } from '@angular/core';
import { AuthService } from './auth/auth.service';
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

  constructor(private characterService: CharacterService, private authService: AuthService) {
    // the AppComponent is loaded regardless wether we are logged in or not,
    // however we can only display certain values if we are logged in
    if (this.authService.isLoggedIn()) {
      this.characterService.getCharacter().subscribe(character => {
        this.character = character;
        this.characterPortrait = this.characterService.getCharacterPortraitURL(character.characterID, 64);
      });
    }
  }

  ngOnInit() {

  }

}
