import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { IonicModule } from '@ionic/angular';
import { Firestore, collection, collectionData} from '@angular/fire/firestore';
import { Observable } from 'rxjs';

export interface Product {
  name: string;
  // (source supplier : price) pairs
  prices: [string: number];
  img_url: string;
  stock: number;
}

@Component({
  selector: 'app-catalogue',
  templateUrl: './catalogue.page.html',
  styleUrls: ['./catalogue.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule, FormsModule]
})
export class CataloguePage implements OnInit {
  private firestore: Firestore = inject(Firestore);
  products$: Observable<Product[]>;

  constructor(store: Firestore) {
    this.products$ = collectionData(collection(store, 'products')) as Observable<Product[]>;
  }

  ngOnInit() {
  }

}
