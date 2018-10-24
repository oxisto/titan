import { Component, OnInit } from '@angular/core';
import 'rxjs/add/operator/map';
import { AuthService } from '../auth/auth.service';
import { ManufacturingService } from '../manufacturing/manufacturing.service';



@Component({
  templateUrl: './blueprints.component.html',
  styleUrls: ['./blueprints.component.css']
})
export class BlueprintsComponent implements OnInit {

  products: any[];

  sortByOptions = ['Profit.PerDay.BasedOnSellPrice:DESC', 'Profit.PerDay.BasedOnBuyPrice:DESC'];

  sortBy: string;

  categories: Map<number, any>;

  selectedCategories: any = {};

  nameFilter: string;

  maxProductionCosts: number;

  constructor(private auth: AuthService,
    private manufacturingService: ManufacturingService) {
    if (!this.auth.isLoggedIn()) {
      return;
    }

    this.maxProductionCosts = +localStorage.getItem('manufacturing:maxProductionCosts');
    this.nameFilter = localStorage.getItem('manufacturing:nameFilter');
    this.sortBy = localStorage.getItem('manufacturing:sortBy');

    if (!this.sortBy) {
      this.sortBy = this.sortByOptions[0];
    }

    this.manufacturingService.getManufacturingCategories().subscribe(categories => {
      this.categories = categories.reduce((map, category) => ({ ...map, [category.categoryID]: category }), {});

      // see, if there is something in localStorage, otherwise set all categories to true
      const json = localStorage.getItem('manufacturing:selectedCategories');
      if (!json) {
        this.categories.forEach(category => {
          this.selectedCategories[category._id] = true;
        });
      } else {
        this.selectedCategories = JSON.parse(json);
      }

      this.fetchProducts();
    });
  }

  ngOnInit() {
  }

  fetchProducts() {
    if (this.auth.isLoggedIn()) {
      const categoryIDs: number[] = [];

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
        maxProductionCosts: this.maxProductionCosts
      }).subscribe(products => {
        this.products = products;
      });
    }
  }

  onSortByChanged(event) {
    // save the selected options in local storage
    localStorage.setItem('manufacturing:maxProductionCosts', this.maxProductionCosts.toString());
    localStorage.setItem('manufacturing:nameFilter', this.nameFilter);
    localStorage.setItem('manufacturing:sortBy', this.sortBy);
    localStorage.setItem('manufacturing:selectedCategories', JSON.stringify(this.selectedCategories));

    this.fetchProducts();
  }

}
