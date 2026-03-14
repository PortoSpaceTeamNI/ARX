export const ValveState = {
  Closed: 0,
  Opened: 1,
  Closing: 2,
  Opening: 3,
  ClosingNotAcked: 4,
  OpeningNotAcked: 5,
} as const;

export type ValveState = (typeof ValveState)[keyof typeof ValveState];

export const ValveID = {
  Pressurizing: 0,
  Vent: 1,
  Abort: 2,
  Main: 3,
  N2OFill: 4,
  N2OPurge: 5,
  N2Fill: 6,
  N2Purge: 7,
  N2OQuickDc: 8,
  N2QuickDc: 9,
} as const;

export type ValveID = (typeof ValveID)[keyof typeof ValveID];

export const ValveNames: Record<ValveID, string> = {
  [ValveID.Pressurizing]: 'Pressurizing',
  [ValveID.Vent]: 'Vent',
  [ValveID.Abort]: 'Abort',
  [ValveID.Main]: 'Main',
  [ValveID.N2OFill]: 'N2O Fill',
  [ValveID.N2OPurge]: 'N2O Purge',
  [ValveID.N2Fill]: 'N2 Fill',
  [ValveID.N2Purge]: 'N2 Purge',
  [ValveID.N2OQuickDc]: 'N2O Quick Disconnect',
  [ValveID.N2QuickDc]: 'N2 Quick Disconnect',
};
