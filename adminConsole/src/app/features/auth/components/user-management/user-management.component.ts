import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AuthService } from '../../../../core/services/auth.service';
import { User, UserRole } from '../../../../core/models/user.model';

@Component({
  selector: 'app-user-management',
  templateUrl: './user-management.component.html',
  styleUrls: ['./user-management.component.scss']
})
export class UserManagementComponent implements OnInit {
  users: User[] = [];
  currentUser: User | null = null;
  showUserModal = false;
  editingUser: User | null = null;
  isSaving = false;
  userForm: FormGroup;
  userRoles = Object.values(UserRole);

  constructor(
    private authService: AuthService,
    private fb: FormBuilder
  ) {
    this.userForm = this.fb.group({
      username: ['', [Validators.required, Validators.minLength(3)]],
      email: ['', [Validators.required, Validators.email]],
      role: [UserRole.VIEWER, [Validators.required]],
      password: ['', [Validators.required, Validators.minLength(6)]],
      isActive: [true]
    });
  }

  ngOnInit(): void {
    this.currentUser = this.authService.getCurrentUser();
    this.loadUsers();
  }

  loadUsers(): void {
    this.authService.getUsers().subscribe(users => {
      this.users = users;
    });
  }

  refreshUsers(): void {
    this.loadUsers();
  }

  openCreateUserModal(): void {
    this.editingUser = null;
    this.userForm.reset({
      role: UserRole.VIEWER,
      isActive: true
    });
    this.showUserModal = true;
  }

  editUser(user: User): void {
    this.editingUser = user;
    this.userForm.patchValue({
      username: user.username,
      email: user.email,
      role: user.role,
      password: '',
      isActive: user.isActive
    });
    this.showUserModal = true;
  }

  viewUserDetails(user: User): void {
    // Implement user details view
    console.log('View user details:', user);
  }

  deleteUser(userId: string): void {
    if (confirm('Are you sure you want to delete this user?')) {
      this.authService.deleteUser(userId).subscribe({
        next: (response) => {
          if (response.success) {
            this.loadUsers();
          }
        },
        error: (error) => {
          console.error('Error deleting user:', error);
        }
      });
    }
  }

  saveUser(): void {
    if (this.userForm.valid) {
      this.isSaving = true;
      const userData = this.userForm.value;

      if (this.editingUser) {
        this.authService.updateUser(this.editingUser.id, userData).subscribe({
          next: (response) => {
            this.isSaving = false;
            if (response.success) {
              this.closeUserModal();
              this.loadUsers();
            }
          },
          error: (error) => {
            this.isSaving = false;
            console.error('Error updating user:', error);
          }
        });
      } else {
        this.authService.createUser(userData).subscribe({
          next: (response) => {
            this.isSaving = false;
            if (response.success) {
              this.closeUserModal();
              this.loadUsers();
            }
          },
          error: (error) => {
            this.isSaving = false;
            console.error('Error creating user:', error);
          }
        });
      }
    }
  }

  closeUserModal(event?: Event): void {
    if (event) {
      const target = event.target as HTMLElement;
      if (target.classList.contains('modal')) {
        this.showUserModal = false;
      }
    } else {
      this.showUserModal = false;
    }
  }

  getRoleClass(role: UserRole): string {
    switch (role) {
      case UserRole.ADMIN:
        return 'role-admin';
      case UserRole.SECURITY_ANALYST:
        return 'role-analyst';
      case UserRole.DEVELOPER:
        return 'role-developer';
      case UserRole.VIEWER:
        return 'role-viewer';
      default:
        return 'role-default';
    }
  }

  getActiveUsersCount(): number {
    return this.users.filter(user => user.isActive).length;
  }

  getUniqueRolesCount(): number {
    const roles = new Set(this.users.map(user => user.role));
    return roles.size;
  }
}
