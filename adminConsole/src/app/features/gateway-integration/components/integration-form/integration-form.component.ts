import { Component, OnInit, OnDestroy } from '@angular/core';
import { FormBuilder, FormGroup, Validators, FormArray } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { GatewayIntegrationService } from '../../services/gateway-integration.service';
import { Integration, GatewayType, IntegrationStatus, Endpoint, Credentials } from '../../../../core/models/gateway-integration.model';
import { Subject, takeUntil } from 'rxjs';

@Component({
  selector: 'app-integration-form',
  templateUrl: './integration-form.component.html',
  styleUrls: ['./integration-form.component.scss']
})
export class IntegrationFormComponent implements OnInit, OnDestroy {
  integrationForm: FormGroup;
  loading = false;
  saving = false;
  error: string | null = null;
  isEditMode = false;
  integrationId: string | null = null;
  private destroy$ = new Subject<void>();

  // Available gateway types
  gatewayTypes = [
    { value: 'kong', label: 'Kong', icon: '🔗', description: 'API Gateway and Microservice Management' },
    { value: 'nginx', label: 'NGINX', icon: '⚡', description: 'High-performance HTTP Server and Load Balancer' },
    { value: 'traefik', label: 'Traefik', icon: '🚦', description: 'Modern HTTP Reverse Proxy and Load Balancer' },
    { value: 'envoy', label: 'Envoy', icon: '🛡️', description: 'High-performance C++ L7 proxy and communication bus' },
    { value: 'haproxy', label: 'HAProxy', icon: '⚖️', description: 'Reliable, High Performance TCP/HTTP Load Balancer' }
  ];

  // Credential types
  credentialTypes = [
    { value: 'basic', label: 'Basic Auth' },
    { value: 'token', label: 'Bearer Token' },
    { value: 'api_key', label: 'API Key' },
    { value: 'tls', label: 'TLS Certificate' }
  ];

  // Protocol options
  protocols = ['http', 'https', 'tcp', 'udp'];

  constructor(
    private fb: FormBuilder,
    private gatewayIntegrationService: GatewayIntegrationService,
    private route: ActivatedRoute,
    private router: Router
  ) {
    this.integrationForm = this.createForm();
  }

  ngOnInit(): void {
    this.route.params.pipe(takeUntil(this.destroy$)).subscribe(params => {
      if (params['id']) {
        this.isEditMode = true;
        this.integrationId = params['id'];
        this.loadIntegration(this.integrationId!);
      }
    });

    // Watch for gateway type changes to update form fields
    this.integrationForm.get('type')?.valueChanges.pipe(takeUntil(this.destroy$)).subscribe(type => {
      this.onGatewayTypeChange(type);
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  createForm(): FormGroup {
    return this.fb.group({
      name: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(100)]],
      type: ['', Validators.required],
      status: ['pending', Validators.required],
      config: this.fb.group({}),
      credentials: this.fb.group({
        type: [''],
        username: [''],
        password: [''],
        token: [''],
        apiKey: [''],
        certificate: [''],
        privateKey: ['']
      }),
      endpoints: this.fb.array([])
    });
  }

  loadIntegration(id: string): void {
    this.loading = true;
    this.error = null;

    this.gatewayIntegrationService.getIntegration(id).subscribe({
      next: (integration) => {
        this.populateForm(integration);
        this.loading = false;
      },
      error: (error) => {
        this.error = 'Failed to load integration: ' + error.message;
        this.loading = false;
      }
    });
  }

  populateForm(integration: Integration): void {
    // Set basic fields
    this.integrationForm.patchValue({
      name: integration.name,
      type: integration.type,
      status: integration.status
    });

    // Set configuration
    if (integration.config) {
      const configGroup = this.integrationForm.get('config') as FormGroup;
      Object.keys(integration.config).forEach(key => {
        if (!configGroup.contains(key)) {
          configGroup.addControl(key, this.fb.control(integration.config[key]));
        } else {
          configGroup.get(key)?.setValue(integration.config[key]);
        }
      });
    }

    // Set credentials
    if (integration.credentials) {
      this.integrationForm.patchValue({
        credentials: {
          type: integration.credentials.type,
          username: integration.credentials.username || '',
          password: integration.credentials.password || '',
          token: integration.credentials.token || '',
          apiKey: integration.credentials.apiKey || '',
          certificate: integration.credentials.certificate || '',
          privateKey: integration.credentials.privateKey || ''
        }
      });
    }

    // Set endpoints
    const endpointsArray = this.integrationForm.get('endpoints') as FormArray;
    endpointsArray.clear();
    integration.endpoints.forEach(endpoint => {
      endpointsArray.push(this.createEndpointFormGroup(endpoint));
    });
  }

  onGatewayTypeChange(type: GatewayType): void {
    const configGroup = this.integrationForm.get('config') as FormGroup;
    
    // Clear existing config controls
    Object.keys(configGroup.controls).forEach(key => {
      configGroup.removeControl(key);
    });

    // Add type-specific configuration fields
    switch (type) {
      case 'kong':
        configGroup.addControl('adminUrl', this.fb.control('', Validators.required));
        configGroup.addControl('proxyUrl', this.fb.control('', Validators.required));
        configGroup.addControl('timeout', this.fb.control(30000));
        break;
      case 'nginx':
        configGroup.addControl('configPath', this.fb.control('/etc/nginx/nginx.conf', Validators.required));
        configGroup.addControl('reloadCommand', this.fb.control('nginx -s reload'));
        configGroup.addControl('backupConfig', this.fb.control(true));
        break;
      case 'traefik':
        configGroup.addControl('apiUrl', this.fb.control('http://localhost:8080', Validators.required));
        configGroup.addControl('dashboard', this.fb.control(true));
        configGroup.addControl('insecureSkipVerify', this.fb.control(false));
        break;
      case 'envoy':
        configGroup.addControl('adminUrl', this.fb.control('http://localhost:9901', Validators.required));
        configGroup.addControl('configPath', this.fb.control('/etc/envoy/envoy.yaml'));
        configGroup.addControl('hotRestart', this.fb.control(false));
        break;
      case 'haproxy':
        configGroup.addControl('configPath', this.fb.control('/etc/haproxy/haproxy.cfg', Validators.required));
        configGroup.addControl('socketPath', this.fb.control('/var/run/haproxy.sock'));
        configGroup.addControl('reloadCommand', this.fb.control('systemctl reload haproxy'));
        break;
    }
  }

  createEndpointFormGroup(endpoint?: Endpoint): FormGroup {
    return this.fb.group({
      id: [endpoint?.id || ''],
      name: [endpoint?.name || '', Validators.required],
      url: [endpoint?.url || '', [Validators.required, Validators.pattern('^https?://.+')]],
      protocol: [endpoint?.protocol || 'http', Validators.required],
      port: [endpoint?.port || 80, [Validators.required, Validators.min(1), Validators.max(65535)]],
      timeout: [endpoint?.timeout || 30000, [Validators.required, Validators.min(1000)]]
    });
  }

  addEndpoint(): void {
    const endpointsArray = this.integrationForm.get('endpoints') as FormArray;
    endpointsArray.push(this.createEndpointFormGroup());
  }

  removeEndpoint(index: number): void {
    const endpointsArray = this.integrationForm.get('endpoints') as FormArray;
    endpointsArray.removeAt(index);
  }

  get endpointsArray(): FormArray {
    return this.integrationForm.get('endpoints') as FormArray;
  }

  onCredentialTypeChange(): void {
    const credentialsGroup = this.integrationForm.get('credentials') as FormGroup;
    const credentialType = credentialsGroup.get('type')?.value;

    // Reset all credential fields
    ['username', 'password', 'token', 'apiKey', 'certificate', 'privateKey'].forEach(field => {
      credentialsGroup.get(field)?.setValue('');
    });

    // Set required validators based on type
    switch (credentialType) {
      case 'basic':
        credentialsGroup.get('username')?.setValidators([Validators.required]);
        credentialsGroup.get('password')?.setValidators([Validators.required]);
        break;
      case 'token':
        credentialsGroup.get('token')?.setValidators([Validators.required]);
        break;
      case 'api_key':
        credentialsGroup.get('apiKey')?.setValidators([Validators.required]);
        break;
      case 'tls':
        credentialsGroup.get('certificate')?.setValidators([Validators.required]);
        credentialsGroup.get('privateKey')?.setValidators([Validators.required]);
        break;
    }

    // Update validators
    ['username', 'password', 'token', 'apiKey', 'certificate', 'privateKey'].forEach(field => {
      credentialsGroup.get(field)?.updateValueAndValidity();
    });
  }

  onSubmit(): void {
    if (this.integrationForm.valid) {
      this.saving = true;
      this.error = null;

      const formValue = this.integrationForm.value;
      
      // Prepare integration object
      const integration: Partial<Integration> = {
        name: formValue.name,
        type: formValue.type,
        status: formValue.status,
        config: formValue.config,
        endpoints: formValue.endpoints
      };

      // Add credentials if type is selected
      if (formValue.credentials.type) {
        integration.credentials = {
          type: formValue.credentials.type,
          username: formValue.credentials.username || undefined,
          password: formValue.credentials.password || undefined,
          token: formValue.credentials.token || undefined,
          apiKey: formValue.credentials.apiKey || undefined,
          certificate: formValue.credentials.certificate || undefined,
          privateKey: formValue.credentials.privateKey || undefined
        };
      }

      const request = this.isEditMode && this.integrationId
        ? this.gatewayIntegrationService.updateIntegration(this.integrationId, integration)
        : this.gatewayIntegrationService.createIntegration(integration);

      request.subscribe({
        next: (result) => {
          this.saving = false;
          this.router.navigate(['/gateway-integration/integrations', result.id]);
        },
        error: (error) => {
          this.error = 'Failed to save integration: ' + error.message;
          this.saving = false;
        }
      });
    } else {
      this.markFormGroupTouched();
    }
  }

  onCancel(): void {
    this.router.navigate(['/gateway-integration/integrations']);
  }

  onTestConnection(): void {
    if (this.integrationForm.valid) {
      const formValue = this.integrationForm.value;
      // Create a temporary integration for testing
      const testIntegration: Partial<Integration> = {
        name: formValue.name,
        type: formValue.type,
        config: formValue.config,
        endpoints: formValue.endpoints
      };

      if (formValue.credentials.type) {
        testIntegration.credentials = {
          type: formValue.credentials.type,
          username: formValue.credentials.username || undefined,
          password: formValue.credentials.password || undefined,
          token: formValue.credentials.token || undefined,
          apiKey: formValue.credentials.apiKey || undefined,
          certificate: formValue.credentials.certificate || undefined,
          privateKey: formValue.credentials.privateKey || undefined
        };
      }

      // For now, just show a success message
      alert('Connection test would be performed here. In a real implementation, this would test the actual gateway connection.');
    } else {
      this.markFormGroupTouched();
    }
  }

  private markFormGroupTouched(): void {
    Object.keys(this.integrationForm.controls).forEach(key => {
      const control = this.integrationForm.get(key);
      control?.markAsTouched();

      if (control instanceof FormGroup) {
        Object.keys(control.controls).forEach(nestedKey => {
          control.get(nestedKey)?.markAsTouched();
        });
      } else if (control instanceof FormArray) {
        control.controls.forEach(arrayControl => {
          if (arrayControl instanceof FormGroup) {
            Object.keys(arrayControl.controls).forEach(nestedKey => {
              arrayControl.get(nestedKey)?.markAsTouched();
            });
          } else {
            arrayControl.markAsTouched();
          }
        });
      }
    });
  }

  getFieldError(fieldName: string): string {
    const field = this.integrationForm.get(fieldName);
    if (field?.errors && field.touched) {
      if (field.errors['required']) return 'This field is required';
      if (field.errors['minlength']) return `Minimum length is ${field.errors['minlength'].requiredLength} characters`;
      if (field.errors['maxlength']) return `Maximum length is ${field.errors['maxlength'].requiredLength} characters`;
      if (field.errors['min']) return `Minimum value is ${field.errors['min'].min}`;
      if (field.errors['max']) return `Maximum value is ${field.errors['max'].max}`;
      if (field.errors['pattern']) return 'Invalid format';
    }
    return '';
  }

  getEndpointFieldError(endpointIndex: number, fieldName: string): string {
    const endpoint = this.endpointsArray.at(endpointIndex);
    const field = endpoint.get(fieldName);
    if (field?.errors && field.touched) {
      if (field.errors['required']) return 'This field is required';
      if (field.errors['min']) return `Minimum value is ${field.errors['min'].min}`;
      if (field.errors['max']) return `Maximum value is ${field.errors['max'].max}`;
      if (field.errors['pattern']) return 'Invalid URL format';
    }
    return '';
  }

  clearError(): void {
    this.error = null;
  }
} 