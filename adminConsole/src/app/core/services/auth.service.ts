import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, of } from 'rxjs';
import { Router } from '@angular/router';
import { User, UserRole, LoginRequest, LoginResponse } from '../models/user.model';

export interface AccessPolicy {
  id: string;
  name: string;
  description: string;
  resources: string[];
  actions: string[];
  conditions: any;
  isActive: boolean;
  createdAt: Date;
}

export interface LoginCredentials {
  username: string;
  password: string;
}

export interface RegisterData {
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private currentUserSubject = new BehaviorSubject<User | null>(null);
  public currentUser$ = this.currentUserSubject.asObservable();

  private isAuthenticatedSubject = new BehaviorSubject<boolean>(false);
  public isAuthenticated$ = this.isAuthenticatedSubject.asObservable();

  constructor(private router: Router) {
    this.loadUserFromStorage();
  }

  // Login functionality
  login(credentials: LoginCredentials): Observable<{ success: boolean; message: string }> {
    // Simulate API call
    return new Observable(observer => {
      setTimeout(() => {
        if (credentials.username === 'admin' && credentials.password === 'admin123') {
          const user: User = {
            id: '1',
            username: 'admin',
            email: 'admin@scopeapi.com',
            firstName: 'Admin',
            lastName: 'User',
            role: UserRole.ADMIN,
            permissions: ['*'],
            isActive: true,
            lastLogin: new Date(),
            createdAt: new Date('2024-01-01'),
            updatedAt: new Date()
          };
          
          this.setCurrentUser(user);
          observer.next({ success: true, message: 'Login successful' });
        } else {
          observer.next({ success: false, message: 'Invalid credentials' });
        }
        observer.complete();
      }, 1000);
    });
  }

  // Logout functionality
  logout(): void {
    localStorage.removeItem('currentUser');
    localStorage.removeItem('token');
    this.currentUserSubject.next(null);
    this.isAuthenticatedSubject.next(false);
    this.router.navigate(['/auth/login']);
  }

  // Get current user
  getCurrentUser(): User | null {
    return this.currentUserSubject.value;
  }

  // Check if user is authenticated
  isAuthenticated(): boolean {
    return this.isAuthenticatedSubject.value;
  }

  // Check if user has specific permission
  hasPermission(permission: string): boolean {
    const user = this.getCurrentUser();
    return user ? user.permissions.includes('*') || user.permissions.includes(permission) : false;
  }

  // Check if user has specific role
  hasRole(roleName: string): boolean {
    const user = this.getCurrentUser();
    return user ? user.role === roleName : false;
  }

  // Get all users (for admin functionality)
  getUsers(): Observable<User[]> {
    // Simulate API call
    return of([
      {
        id: '1',
        username: 'admin',
        email: 'admin@scopeapi.com',
        firstName: 'Admin',
        lastName: 'User',
        role: UserRole.ADMIN,
        permissions: ['*'],
        isActive: true,
        lastLogin: new Date(),
        createdAt: new Date('2024-01-01'),
        updatedAt: new Date()
      },
      {
        id: '2',
        username: 'analyst',
        email: 'analyst@scopeapi.com',
        firstName: 'Security',
        lastName: 'Analyst',
        role: UserRole.SECURITY_ANALYST,
        permissions: ['threat:view', 'threat:manage'],
        isActive: true,
        lastLogin: new Date(),
        createdAt: new Date('2024-01-01'),
        updatedAt: new Date()
      }
    ]);
  }

  // Create new user
  createUser(userData: any): Observable<{ success: boolean; message: string }> {
    // Simulate API call
    return of({ success: true, message: 'User created successfully' });
  }

  // Update existing user
  updateUser(userId: string, userData: any): Observable<{ success: boolean; message: string }> {
    // Simulate API call
    return of({ success: true, message: 'User updated successfully' });
  }

  // Delete user
  deleteUser(userId: string): Observable<{ success: boolean; message: string }> {
    // Simulate API call
    return of({ success: true, message: 'User deleted successfully' });
  }

  // Get access policies
  getAccessPolicies(): Observable<AccessPolicy[]> {
    // Simulate API call
    return of([
      {
        id: '1',
        name: 'Admin Full Access',
        description: 'Full system access for administrators',
        resources: ['*'],
        actions: ['*'],
        conditions: {},
        isActive: true,
        createdAt: new Date('2024-01-01')
      },
      {
        id: '2',
        name: 'Analyst Read Access',
        description: 'Read-only access for security analysts',
        resources: ['threats', 'endpoints'],
        actions: ['read'],
        conditions: {},
        isActive: true,
        createdAt: new Date('2024-01-01')
      }
    ]);
  }

  // Create access policy
  createAccessPolicy(policyData: any): Observable<{ success: boolean; message: string }> {
    // Simulate API call
    return of({ success: true, message: 'Access policy created successfully' });
  }

  // Update access policy
  updateAccessPolicy(policyId: string, policyData: any): Observable<{ success: boolean; message: string }> {
    // Simulate API call
    return of({ success: true, message: 'Access policy updated successfully' });
  }

  // Delete access policy
  deleteAccessPolicy(policyId: string): Observable<{ success: boolean; message: string }> {
    // Simulate API call
    return of({ success: true, message: 'Access policy deleted successfully' });
  }

  // Private methods
  private setCurrentUser(user: User): void {
    localStorage.setItem('currentUser', JSON.stringify(user));
    this.currentUserSubject.next(user);
    this.isAuthenticatedSubject.next(true);
  }

  private loadUserFromStorage(): void {
    const storedUser = localStorage.getItem('currentUser');
    if (storedUser) {
      const user = JSON.parse(storedUser);
      this.currentUserSubject.next(user);
      this.isAuthenticatedSubject.next(true);
    }
  }
}
