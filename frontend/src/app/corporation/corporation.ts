import { Character } from '../character/character';

export class Corporation {

  name: string;

  corporationID: number;

  allianceID: number;

  CEOID: number;

  ticker: string;

  members: Map<Number, Character> = new Map();

  expiresOn: Date;

  getId(): number {
    return this.corporationID;
  }

}
