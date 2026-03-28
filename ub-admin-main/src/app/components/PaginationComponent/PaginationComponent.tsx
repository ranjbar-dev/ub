import { Pagination } from '@material-ui/lab';
import React, { memo } from 'react';
/**
 * Pagination bar rendered below a SimpleGrid.
 * Wraps Material-UI Pagination with first/last page buttons.
 *
 * @example
 * ```tsx
 * <PaginationComponent size={10} onPageChange={(page) => loadPage(page)} />
 * ```
 */
function PaginationComponent(props: { onPageChange: (page: number) => void; size: number; style?: React.CSSProperties }) {
	const {style}=props
  return (
    <div>
      <Pagination
        style={{ position: 'absolute', right: '5px',...style }}
        count={props.size}
        onChange={(e, page) => {
          props.onPageChange(page);
        }}
        showFirstButton
        showLastButton
        variant="outlined"
        //  shape="rounded"
      />
    </div>
  );
}

export default memo(PaginationComponent);
