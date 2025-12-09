import { Routes } from '@angular/router';
import { FileList } from './components/file-list/file-list'; 

export const routes: Routes = [  {
    path: '',
    component: FileList,
    title: 'Files',
  }];
