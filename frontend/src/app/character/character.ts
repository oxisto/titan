export class Character {
  name: string;
  characterID: number;

  corporationID: number;
  corporationName: string;

  allianceID: number;
  allianceName: string;

  skills: Map<Number, Skill> = new Map();

  expiresOn: Date;

  getId(): number {
    return this.characterID;
  }
}

export class Skill {
  skillID: number;
  skillpoints: number;
  level: number;

  constructor({ skillID, level, skillpoints }) {
    this.skillID = skillID;
    this.level = level;
    this.skillpoints = skillpoints;
  }

  getId() {
    return this.skillID;
  }
}
