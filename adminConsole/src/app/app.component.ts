import { Component } from '@angular/core';
import { Router, NavigationEnd } from '@angular/router';
import { filter } from 'rxjs/operators';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'admin-console';
  currentUrl: string = '';

  constructor(private router: Router) {
    this.router.events.pipe(
      filter(event => event instanceof NavigationEnd)
    ).subscribe((event: any) => {
      if (event instanceof NavigationEnd) {
        this.currentUrl = event.url;
      }
    });
  }

  isAuthPage(): boolean {
    return (
      this.currentUrl.startsWith('/auth/login') ||
      this.currentUrl.startsWith('/auth/register') ||
      this.currentUrl.startsWith('/auth/forgot-password')
    );
  }
}
