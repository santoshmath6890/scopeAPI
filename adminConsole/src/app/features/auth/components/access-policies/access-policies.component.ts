import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AuthService } from '../../../../core/services/auth.service';
import { AccessPolicy } from '../../../../core/services/auth.service';

@Component({
  selector: 'app-access-policies',
  templateUrl: './access-policies.component.html',
  styleUrls: ['./access-policies.component.scss']
})
export class AccessPoliciesComponent implements OnInit {
  policies: AccessPolicy[] = [];
  showPolicyModal = false;
  editingPolicy: AccessPolicy | null = null;
  isSaving = false;
  policyForm: FormGroup;

  constructor(
    private authService: AuthService,
    private fb: FormBuilder
  ) {
    this.policyForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(3)]],
      description: ['', [Validators.required]],
      resources: [[], [Validators.required]],
      actions: [[], [Validators.required]],
      isActive: [true]
    });
  }

  ngOnInit(): void {
    this.loadPolicies();
  }

  loadPolicies(): void {
    this.authService.getAccessPolicies().subscribe(policies => {
      this.policies = policies;
    });
  }

  refreshPolicies(): void {
    this.loadPolicies();
  }

  openCreatePolicyModal(): void {
    this.editingPolicy = null;
    this.policyForm.reset({
      isActive: true,
      resources: [],
      actions: []
    });
    this.showPolicyModal = true;
  }

  editPolicy(policy: AccessPolicy): void {
    this.editingPolicy = policy;
    this.policyForm.patchValue({
      name: policy.name,
      description: policy.description,
      resources: policy.resources,
      actions: policy.actions,
      isActive: policy.isActive
    });
    this.showPolicyModal = true;
  }

  viewPolicyDetails(policy: AccessPolicy): void {
    // Implement policy details view
    console.log('View policy details:', policy);
  }

  deletePolicy(policyId: string): void {
    if (confirm('Are you sure you want to delete this policy?')) {
      this.authService.deleteAccessPolicy(policyId).subscribe({
        next: (response) => {
          if (response.success) {
            this.loadPolicies();
          }
        },
        error: (error) => {
          console.error('Error deleting policy:', error);
        }
      });
    }
  }

  savePolicy(): void {
    if (this.policyForm.valid) {
      this.isSaving = true;
      const policyData = this.policyForm.value;

      if (this.editingPolicy) {
        this.authService.updateAccessPolicy(this.editingPolicy.id, policyData).subscribe({
          next: (response) => {
            this.isSaving = false;
            if (response.success) {
              this.closePolicyModal();
              this.loadPolicies();
            }
          },
          error: (error) => {
            this.isSaving = false;
            console.error('Error updating policy:', error);
          }
        });
      } else {
        this.authService.createAccessPolicy(policyData).subscribe({
          next: (response) => {
            this.isSaving = false;
            if (response.success) {
              this.closePolicyModal();
              this.loadPolicies();
            }
          },
          error: (error) => {
            this.isSaving = false;
            console.error('Error creating policy:', error);
          }
        });
      }
    }
  }

  closePolicyModal(event?: Event): void {
    if (event) {
      const target = event.target as HTMLElement;
      if (target.classList.contains('modal')) {
        this.showPolicyModal = false;
      }
    } else {
      this.showPolicyModal = false;
    }
  }

  hasConditions(conditions: any): boolean {
    return conditions && Object.keys(conditions).length > 0;
  }

  getConditionsList(conditions: any): Array<{key: string, value: string}> {
    return Object.entries(conditions).map(([key, value]) => ({
      key,
      value: Array.isArray(value) ? value.join(', ') : String(value)
    }));
  }

  getActivePoliciesCount(): number {
    return this.policies.filter(policy => policy.isActive).length;
  }

  getUniqueResourcesCount(): number {
    const resources = new Set();
    this.policies.forEach(policy => {
      policy.resources.forEach(resource => resources.add(resource));
    });
    return resources.size;
  }
}
