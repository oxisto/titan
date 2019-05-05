import { Component, OnInit } from '@angular/core';
import { AuthService } from '../auth/auth.service';
import { ManufacturingService } from '../manufacturing/manufacturing.service';

@Component({
  templateUrl: './blueprints.component.html',
  styleUrls: ['./blueprints.component.css']
})
export class BlueprintsComponent implements OnInit {

  products: any[];

  sortByOptions = ['basedOnSellPrice', 'basedOnBuyPrice'];

  sortBy: string;

  categories: Map<number, any>;

  // preferably this would be a Map<number, any>. However angular seems to have issues with two-way bindings of Maps
  selectedCategories: any = {};

  nameFilter: string;
  hasRequiredSkillsOnly: boolean;
  maxProductionCosts: number;

  constructor(private auth: AuthService,
    private manufacturingService: ManufacturingService) {
    this.maxProductionCosts = +localStorage.getItem('manufacturing:maxProductionCosts');
    this.nameFilter = localStorage.getItem('manufacturing:nameFilter');

    if (!this.nameFilter) {
      this.nameFilter = '';
    }

    this.hasRequiredSkillsOnly = localStorage.getItem('manufacturing:hasRequiredSkillsOnly') === 'true';
    this.sortBy = localStorage.getItem('manufacturing:sortBy');

    if (!this.sortBy) {
      this.sortBy = this.sortByOptions[0];
    }

    this.manufacturingService.getManufacturingCategories().subscribe(categories => {
      this.categories = categories.reduce((map, category) => ({ ...map, [category.categoryID]: category }), {});

      // see, if there is something in localStorage, otherwise set all categories to true
      const json = localStorage.getItem('manufacturing:selectedCategories');
      if (!json) {
        this.selectAllCategories();
      } else {
        this.selectedCategories = JSON.parse(json);
      }

      this.fetchProducts();
    });
  }

  ngOnInit() {
  }

  fetchProducts() {
    const categoryIDs: number[] = [];

    console.log(this.selectedCategories);

    for (const key of Object.keys(this.selectedCategories)) {
      const value = this.selectedCategories[key];
      if (value) {
        categoryIDs.push(+key);
      }
    }

    this.manufacturingService.getManufacturingTypeIDs({
      categoryIDs: categoryIDs,
      sortBy: this.sortBy,
      nameFilter: this.nameFilter,
      hasRequiredSkillsOnly: this.hasRequiredSkillsOnly,
      maxProductionCosts: this.maxProductionCosts
    }).subscribe(products => {
      this.products = products;
    });
  }

  onSortByChanged(event) {
    // save the selected options in local storage
    localStorage.setItem('manufacturing:maxProductionCosts', this.maxProductionCosts.toString());
    localStorage.setItem('manufacturing:nameFilter', this.nameFilter);
    localStorage.setItem('manufacturing:hasRequiredSkillsOnly', String(this.hasRequiredSkillsOnly));
    localStorage.setItem('manufacturing:sortBy', this.sortBy);
    localStorage.setItem('manufacturing:selectedCategories', JSON.stringify(this.selectedCategories));

    this.fetchProducts();
  }

  selectAllCategories() {
    Object.values(this.categories).forEach(category => {
      this.selectedCategories[category.categoryID] = true;
    });
  }

  onSelectAll() {
    this.selectAllCategories();

    this.onSortByChanged(null);
  }

  onDeselectAll() {
    Object.values(this.categories).forEach(category => {
      this.selectedCategories[category.categoryID] = false;
    });

    this.onSortByChanged(null);
  }

}
