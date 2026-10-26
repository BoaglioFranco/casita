import type { ComponentProps } from 'react';
import type { ListIconId } from './icon-catalog';
import { getListIconDefinition } from './icon-catalog';

type ListIconProps = ComponentProps<'svg'> & {
  iconId: ListIconId;
  size?: number;
};

export function ListIcon({ iconId, ...props }: ListIconProps) {
  const Icon = getListIconDefinition(iconId).icon;
  return <Icon aria-hidden {...props} />;
}
