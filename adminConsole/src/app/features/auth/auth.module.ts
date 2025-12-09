import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';

import { AuthRoutingModule } from './auth-routing.module';
import { LoginComponent } from './components/login/login.component';
import { UserManagementComponent } from './components/user-management/user-management.component';
import { AccessPoliciesComponent } from './components/access-policies/access-policies.component';

@NgModule({
  declarations: [
    LoginComponent,
    UserManagementComponent,
    AccessPoliciesComponent
  ],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule,
    AuthRoutingModule
  ]
})
export class AuthModule { }
