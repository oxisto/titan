import { Character } from '../character/character';

export class Corporation {

  name: string;

  corporationID: string;

  allianceID: string;

  CEOID: string;

  ticker: string;

  members: Map<Number, Character> = new Map();

  expiresOn: Date;

  getId(): string {
    return this.corporationID;
  }

}
