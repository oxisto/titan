import { Component, OnInit } from '@angular/core';
import { Headers, Http } from '@angular/http';

import 'rxjs/add/operator/map';

import { AuthService } from '../auth/auth.service';
import { CorporationService } from '../corporation/corporation.service';
import { CharacterService } from '../character/character.service';
import { Character } from '../character/character';

@Component({
  selector: 'app-sidebar',
  templateUrl: './sidebar.component.html',
  styleUrls: ['./sidebar.component.css']
})
export class SidebarComponent implements OnInit {

  character: Character;
  characterPortrait: string;
  corporationLogo: string;

  constructor(private auth: AuthService,
              private corporationService: CorporationService,
              private characterService: CharacterService) {

  }

  ngOnInit() {
    if (this.auth.isLoggedIn()) {
      this.characterService.character.subscribe(character => {
        this.character = character;
        this.characterPortrait = this.characterService.getCharacterPortraitURL(character.characterID);
        this.corporationLogo = this.corporationService.getCorporationLogo(character.corporationID);
      });
    }
  }

}
