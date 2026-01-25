<script setup lang="ts">
import Card from 'primevue/card'
import ProgressBar from 'primevue/progressbar'
import Avatar from 'primevue/avatar'

const kpiCards = [
  { 
    title: "Campaign Alerts", 
    value: "123k", 
    timestamp: "11:08, 17 Aug",
    icon: "pi-envelope",
    iconColor: "var(--p-blue-500)",
    iconBg: "var(--p-blue-50)"
  },
  { 
    title: "System Alerts", 
    value: "89k", 
    timestamp: "14:15, 18 Aug",
    icon: "pi-bolt",
    iconColor: "var(--p-orange-500)",
    iconBg: "var(--p-orange-50)"
  },
  { 
    title: "Promotional Offers", 
    value: "3k", 
    timestamp: "09:30, 16 Aug",
    icon: "pi-gift",
    iconColor: "var(--p-gray-500)",
    iconBg: "var(--p-gray-50)"
  },
  { 
    title: "Traffic Distribution", 
    value: "175k", 
    timestamp: "16:45, 19 Aug",
    icon: "pi-chart-line",
    iconColor: "var(--p-purple-500)",
    iconBg: "var(--p-purple-50)"
  },
];

const campaignPerformance = [
  { name: "All Traffic", visits: "1250 Visits", icon: "pi-refresh", iconColor: "var(--p-green-500)" },
  { name: "Instagram", visits: "660 Visits", icon: "pi-instagram", iconColor: "var(--p-pink-500)" },
  { name: "Google", visits: "817 Visits", icon: "pi-google", iconColor: "var(--p-blue-500)" },
  { name: "Linkedin", visits: "733 Visits", icon: "pi-linkedin", iconColor: "var(--p-blue-600)" },
  { name: "X", visits: "995 Visits", icon: "pi-twitter", iconColor: "var(--p-gray-900)" },
];

const campaignTargets = {
  overall: { percentage: 85.7, value: "8571/10000" },
  details: [
    { name: "New Subscriptions", current: 152, target: 300, color: "var(--p-blue-500)" },
    { name: "Renewal Contracts", current: 63, target: 500, color: "var(--p-orange-500)" },
    { name: "Upsell Revenue", current: 23, target: 1000, color: "var(--p-purple-500)" },
    { name: "Add-On Sales", current: 42, target: 2000, color: "var(--p-pink-500)" },
  ]
};

const teamPerformance = [
  {
    name: "Cameron Williamson",
    title: "Marketing Coordinator",
    avatar: "CW",
    metrics: [
      { platform: "Twitter", percentage: 34.00 },
      { platform: "Facebook", percentage: 45.86 },
      { platform: "Google", percentage: 79.00 },
    ]
  },
  {
    name: "Kathryn Murphy",
    title: "President of Sales",
    avatar: "KM",
    metrics: [
      { platform: "Twitter", percentage: 64.47 },
      { platform: "Facebook", percentage: 75.67 },
      { platform: "Google", percentage: 45.00 },
    ]
  },
  {
    name: "Darrell Steward",
    title: "Web Designer",
    avatar: "DS",
    metrics: [
      { platform: "Twitter", percentage: 23.55 },
      { platform: "Facebook", percentage: 78.65 },
      { platform: "Google", percentage: 86.54 },
    ]
  },
];
</script>

<template>
  <div class="page-root">
    <div class="page-header">
      <h1>Dashboard</h1>
      <p class="text-muted-foreground">Overview of your system and recent activity</p>
    </div>

    <!-- KPI Cards -->
    <div class="kpi-grid">
      <Card v-for="kpi in kpiCards" :key="kpi.title" class="kpi-card">
        <template #content>
          <div class="kpi-card-content">
            <div class="kpi-card-header">
              <div class="kpi-icon" :style="{ backgroundColor: kpi.iconBg, color: kpi.iconColor }">
                <i :class="'pi ' + kpi.icon"></i>
              </div>
              <span class="kpi-timestamp">{{ kpi.timestamp }}</span>
            </div>
            <h3 class="kpi-title">{{ kpi.title }}</h3>
            <div class="kpi-value">{{ kpi.value }}</div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Middle Row: Campaign Performance & Targets -->
    <div class="dashboard-grid">
      <!-- Campaign Performance -->
      <Card class="campaign-card">
        <template #title>Campaign Performance</template>
        <template #content>
          <div class="campaign-list">
            <div v-for="item in campaignPerformance" :key="item.name" class="campaign-item">
              <div class="campaign-item-left">
                <div class="campaign-icon" :style="{ backgroundColor: item.iconColor + '20', color: item.iconColor }">
                  <i :class="'pi ' + item.icon"></i>
                </div>
                <div>
                  <div class="campaign-name">{{ item.name }}</div>
                  <div class="campaign-visits">{{ item.visits }}</div>
                </div>
              </div>
              <div class="campaign-item-right">
                <div class="campaign-sparkline"></div>
                <i class="pi pi-angle-right text-muted-foreground"></i>
              </div>
            </div>
          </div>
        </template>
      </Card>

      <!-- Campaign Targets -->
      <Card class="targets-card">
        <template #title>Campaign Targets</template>
        <template #content>
          <div class="targets-content">
            <div class="targets-overall">
              <div class="targets-percentage">{{ campaignTargets.overall.percentage }}%</div>
              <div class="targets-value">{{ campaignTargets.overall.value }}</div>
            </div>
            <div class="targets-progress">
              <div 
                v-for="(detail, idx) in campaignTargets.details" 
                :key="detail.name"
                class="targets-segment"
                :style="{ 
                  width: `${(detail.current / detail.target) * 100}%`,
                  backgroundColor: detail.color
                }"
              ></div>
            </div>
            <div class="targets-details">
              <div v-for="detail in campaignTargets.details" :key="detail.name" class="target-detail-item">
                <div class="target-dot" :style="{ backgroundColor: detail.color }"></div>
                <span class="target-detail-name">{{ detail.name }}</span>
                <span class="target-detail-value">{{ detail.current }} / {{ detail.target }}</span>
              </div>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Team Performance Cards -->
    <div class="team-grid">
      <Card v-for="member in teamPerformance" :key="member.name" class="team-card">
        <template #content>
          <div class="team-card-content">
            <Avatar :label="member.avatar" shape="circle" size="large" class="team-avatar" />
            <div class="team-info">
              <h3 class="team-name">{{ member.name }}</h3>
              <p class="team-title">{{ member.title }}</p>
            </div>
            <div class="team-metrics">
              <div v-for="metric in member.metrics" :key="metric.platform" class="team-metric">
                <div class="team-metric-header">
                  <span class="team-metric-platform">{{ metric.platform }}</span>
                  <span class="team-metric-value">{{ metric.percentage.toFixed(2) }}%</span>
                </div>
                <ProgressBar :value="metric.percentage" :show-value="false" class="team-progress" />
              </div>
            </div>
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>

<style scoped>
.page-root { display: flex; flex-direction: column; gap: 1.5rem; }
.page-header h1 { font-size: 1.5rem; font-weight: 700; margin-bottom: 0.25rem; }
.page-header p { font-size: 0.875rem; }

/* KPI Cards */
.kpi-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1rem;
}
.kpi-card {
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}
.kpi-card-content { padding: 1.25rem; }
.kpi-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}
.kpi-icon {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.25rem;
}
.kpi-timestamp {
  font-size: 0.75rem;
  color: var(--p-text-muted-color);
}
.kpi-title {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--p-text-muted-color);
  margin-bottom: 0.5rem;
}
.kpi-value {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--p-text-color);
}

/* Dashboard Grid */
.dashboard-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 1rem;
}
@media (max-width: 1024px) {
  .dashboard-grid { grid-template-columns: 1fr; }
}

/* Campaign Performance */
.campaign-card {
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}
.campaign-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}
.campaign-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem;
  border-radius: var(--p-border-radius);
  transition: background-color 0.2s;
}
.campaign-item:hover {
  background: var(--p-surface-50);
}
.campaign-item-left {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}
.campaign-icon {
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.875rem;
}
.campaign-name {
  font-weight: 500;
  font-size: 0.875rem;
  margin-bottom: 0.25rem;
}
.campaign-visits {
  font-size: 0.75rem;
  color: var(--p-text-muted-color);
}
.campaign-item-right {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
.campaign-sparkline {
  width: 3rem;
  height: 1.5rem;
  background: linear-gradient(to right, var(--p-green-500), var(--p-green-300));
  border-radius: 2px;
}

/* Campaign Targets */
.targets-card {
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}
.targets-content {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.targets-overall {
  display: flex;
  align-items: baseline;
  gap: 0.75rem;
}
.targets-percentage {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--p-text-color);
}
.targets-value {
  font-size: 1rem;
  color: var(--p-text-muted-color);
}
.targets-progress {
  height: 0.5rem;
  background: var(--p-surface-200);
  border-radius: var(--p-border-radius);
  display: flex;
  overflow: hidden;
}
.targets-segment {
  height: 100%;
  transition: width 0.3s;
}
.targets-details {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}
.target-detail-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
.target-dot {
  width: 0.5rem;
  height: 0.5rem;
  border-radius: 50%;
  flex-shrink: 0;
}
.target-detail-name {
  flex: 1;
  font-size: 0.875rem;
  color: var(--p-text-color);
}
.target-detail-value {
  font-size: 0.875rem;
  color: var(--p-text-muted-color);
}

/* Team Performance */
.team-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1rem;
}
.team-card {
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}
.team-card-content {
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.team-avatar {
  align-self: flex-start;
}
.team-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}
.team-name {
  font-size: 1rem;
  font-weight: 600;
  margin: 0;
}
.team-title {
  font-size: 0.875rem;
  color: var(--p-text-muted-color);
  margin: 0;
}
.team-metrics {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}
.team-metric {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.team-metric-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.team-metric-platform {
  font-size: 0.875rem;
  color: var(--p-text-color);
}
.team-metric-value {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--p-text-color);
}
.team-progress {
  height: 0.5rem;
}
</style>
