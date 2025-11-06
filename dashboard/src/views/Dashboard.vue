<script setup lang="ts">
import { Icon } from "@iconify/vue";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

// Mock data for statistics
const stats = [
  {
    title: "Total Nodes",
    value: "12",
    change: "+2 from last month",
    changeType: "positive",
    icon: "lucide:server",
  },
  {
    title: "Active Nodes",
    value: "8",
    change: "+3 from last month",
    changeType: "positive",
    icon: "lucide:activity",
  },
  {
    title: "Applications",
    value: "24",
    change: "+5 from last month",
    changeType: "positive",
    icon: "lucide:layers",
  },
  {
    title: "Total Requests",
    value: "12,234",
    change: "+19% from last month",
    changeType: "positive",
    icon: "lucide:trending-up",
  },
];

// Mock data for recent activity
const recentActivity = [
  {
    id: 1,
    user: "Olivia Martin",
    email: "olivia.martin@email.com",
    action: "Deployed application",
    amount: "+$1,999.00",
    time: "2 hours ago",
  },
  {
    id: 2,
    user: "Jackson Lee",
    email: "jackson.lee@email.com",
    action: "Created node",
    amount: "+$39.00",
    time: "3 hours ago",
  },
  {
    id: 3,
    user: "Isabella Nguyen",
    email: "isabella.nguyen@email.com",
    action: "Updated configuration",
    amount: "+$299.00",
    time: "5 hours ago",
  },
  {
    id: 4,
    user: "William Kim",
    email: "william.kim@email.com",
    action: "Restarted service",
    amount: "+$99.00",
    time: "1 day ago",
  },
  {
    id: 5,
    user: "Sofia Davis",
    email: "sofia.davis@email.com",
    action: "Added new node",
    amount: "+$39.00",
    time: "2 days ago",
  },
];
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-3xl font-bold">Dashboard</h1>
      <p class="text-muted-foreground mt-1">
        Overview of your system and recent activity
      </p>
    </div>

    <!-- Statistics Cards -->
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <Card v-for="stat in stats" :key="stat.title">
        <CardHeader
          class="flex flex-row items-center justify-between space-y-0 pb-2"
        >
          <CardTitle class="text-sm font-medium">
            {{ stat.title }}
          </CardTitle>
          <Icon :icon="stat.icon" class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stat.value }}</div>
          <p class="text-xs text-muted-foreground mt-1">
            <span
              :class="
                stat.changeType === 'positive'
                  ? 'text-green-500'
                  : 'text-red-500'
              "
            >
              {{ stat.change }}
            </span>
          </p>
        </CardContent>
      </Card>
    </div>

    <!-- Overview and Recent Activity -->
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
      <!-- Overview Card -->
      <Card class="col-span-4">
        <CardHeader>
          <CardTitle>Overview</CardTitle>
          <CardDescription>
            System performance and activity overview
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div
            class="flex items-center justify-center h-[300px] text-muted-foreground"
          >
            <div class="text-center">
              <Icon
                icon="lucide:bar-chart-3"
                class="h-12 w-12 mx-auto mb-4 opacity-50"
              />
              <p>Chart visualization will be added here</p>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Recent Activity Card -->
      <Card class="col-span-3">
        <CardHeader>
          <CardTitle>Recent Activity</CardTitle>
          <CardDescription>
            You made 265 operations this month.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div
              v-for="activity in recentActivity"
              :key="activity.id"
              class="flex items-center gap-4"
            >
              <div
                class="flex h-9 w-9 items-center justify-center rounded-full bg-primary/10 text-primary"
              >
                <Icon icon="lucide:user" class="h-4 w-4" />
              </div>
              <div class="flex-1 space-y-1">
                <p class="text-sm font-medium leading-none">
                  {{ activity.user }}
                </p>
                <p class="text-xs text-muted-foreground">
                  {{ activity.email }}
                </p>
                <p class="text-xs text-muted-foreground">
                  {{ activity.action }}
                </p>
              </div>
              <div class="text-right">
                <p class="text-sm font-medium">{{ activity.amount }}</p>
                <p class="text-xs text-muted-foreground">{{ activity.time }}</p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Recent Sales Table -->
    <Card>
      <CardHeader>
        <CardTitle>Recent Operations</CardTitle>
        <CardDescription>
          A list of recent operations and activities in your system.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>User</TableHead>
              <TableHead>Email</TableHead>
              <TableHead>Action</TableHead>
              <TableHead>Amount</TableHead>
              <TableHead>Time</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="activity in recentActivity" :key="activity.id">
              <TableCell class="font-medium">{{ activity.user }}</TableCell>
              <TableCell>{{ activity.email }}</TableCell>
              <TableCell>{{ activity.action }}</TableCell>
              <TableCell>{{ activity.amount }}</TableCell>
              <TableCell class="text-muted-foreground">{{
                activity.time
              }}</TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  </div>
</template>
