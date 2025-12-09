import { Component, OnInit } from '@angular/core';
import { ThreatService } from '../../../../core/services/threat.service';
import { Threat, ThreatStatistics } from '../../../../core/models/threat.model';

@Component({
  selector: 'app-threat-detection-overview',
  templateUrl: './threat-detection-overview.component.html',
  styleUrls: ['./threat-detection-overview.component.scss']
})
export class ThreatDetectionOverviewComponent implements OnInit {
  threats: Threat[] = [];
  threatStatistics: ThreatStatistics | null = null;
  isLoading = false;

  constructor(private threatService: ThreatService) {}

  ngOnInit(): void {
    this.loadThreats();
    this.loadThreatStatistics();
  }

  loadThreats(): void {
    this.isLoading = true;
    this.threatService.getThreats().subscribe({
      next: (threats) => {
        this.threats = threats;
        this.isLoading = false;
      },
      error: (error) => {
        console.error('Error loading threats:', error);
        this.isLoading = false;
      }
    });
  }

  loadThreatStatistics(): void {
    this.threatService.getThreatStatistics().subscribe({
      next: (statistics) => {
        this.threatStatistics = statistics;
      },
      error: (error) => {
        console.error('Error loading threat statistics:', error);
      }
    });
  }

  updateThreatStatus(threatId: string, status: string): void {
    this.threatService.updateThreatStatus(threatId, status as any).subscribe({
      next: (result) => {
        if (result.success) {
          console.log('Threat status updated:', result.message);
          this.loadThreats(); // Refresh the list
        }
      },
      error: (error) => {
        console.error('Error updating threat status:', error);
      }
    });
  }

  blockThreat(threatId: string): void {
    this.threatService.blockThreat(threatId).subscribe({
      next: (result) => {
        if (result.success) {
          console.log('Threat blocked:', result.message);
          this.loadThreats(); // Refresh the list
        }
      },
      error: (error) => {
        console.error('Error blocking threat:', error);
      }
    });
  }
}
