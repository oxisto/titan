import { Character } from '../character/character';

export class Corporation {

  name: string;

  corporationID: number;

  allianceID: number;

  CEOID: number;

  ticker: string;

  members: Map<number, Character> = new Map();

  expiresOn: Date;

  getId(): number {
    return this.corporationID;
  }

}

export class Wallets {

  corporationID: number;
  divisions: Array<Wallet>;

}

export class Wallet {
  balance: number;
  division: number;
}
