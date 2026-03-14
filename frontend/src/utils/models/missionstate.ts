export const MissionState = {
  Idle: 0,
  Filling: 1,
  SafeIdle: 2,
  FillingN2: 3,
  PrePressure: 4,
  FillingN2O: 5,
  PostPressure: 6,
  Ready: 7,
  Armed: 8,
  Ignition: 9,
  Launch: 10,
  Flight: 11,
  Recovery: 12,
  Abort: 13,
} as const;

export type MissionState = (typeof MissionState)[keyof typeof MissionState];

export const MissionStateNames: Record<MissionState, string> = {
  [MissionState.Idle]: 'Idle',
  [MissionState.Filling]: 'Filling',
  [MissionState.SafeIdle]: 'Safe Idle',
  [MissionState.FillingN2]: 'Fill N2',
  [MissionState.PrePressure]: 'Pre Pressurize',
  [MissionState.FillingN2O]: 'Fill N2O',
  [MissionState.PostPressure]: 'Post Pressurize',
  [MissionState.Ready]: 'Ready',
  [MissionState.Armed]: 'Armed',
  [MissionState.Ignition]: 'Ignition',
  [MissionState.Launch]: 'Launch',
  [MissionState.Flight]: 'Flight',
  [MissionState.Recovery]: 'Recovery',
  [MissionState.Abort]: 'Abort',
};
