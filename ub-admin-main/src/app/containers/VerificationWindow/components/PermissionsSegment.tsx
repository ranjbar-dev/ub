import {
  FormControl,
  FormLabel,
  FormGroup,
} from '@material-ui/core';
import SaveOutlinedIcon from '@material-ui/icons/SaveOutlined';
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import IsLoadingWithTextAuto from 'app/components/isLoadingWithText/isLoadingWithTextAuto';
import UbCheckbox from 'app/components/UbCheckBox/UbCheckbox';
import { Buttons } from 'app/constants';
import { InitialUserDetails } from 'app/containers/UserAccounts/types';
import React, { memo, useEffect, useState, useRef } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import styled from 'styled-components/macro';

import { Permission } from '../../UserDetails/types';
import { selectPermissionsData } from '../selectors';
import { VerificationWindowActions } from '../slice';

interface Props {
  data: InitialUserDetails;
}

function PermissionsSegment(props: Props) {
  const [IsLoading, setIsLoading] = useState(true);
  const permissionsArray = useRef<number[]>([]);
  const [Permissions, setPermissions] = useState<Permission[]>([]);
  const { data } = props;
  const dispatch = useDispatch();
  const permissionsState = useSelector(selectPermissionsData);

  useEffect(() => {
    dispatch(VerificationWindowActions.GetPermissionsAction({ id: data.id }));
  }, []);

  useEffect(() => {
    if (
      permissionsState.userId === data.id &&
      permissionsState.data !== null
    ) {
      permissionsArray.current = [];
      for (const item of permissionsState.data) {
        if (item.userHasIt) {
          permissionsArray.current.push(item.id);
        }
      }
      setIsLoading(false);
      setPermissions(permissionsState.data);
    }
  }, [permissionsState]);
  const handleChange = (index: number, checked: boolean, id: number) => {
    let perms: Permission[] = [...Permissions];
    perms[index].userHasIt = checked;
    setPermissions([...perms]);
    if (checked === true) {
      if (!permissionsArray.current.includes(id)) {
        permissionsArray.current.push(id);
      }
    } else {
      const removeIndex = permissionsArray.current.indexOf(id);
      permissionsArray.current.splice(removeIndex, 1);
    }
  };
  const handlePermissionSubmitClick = () => {
    //console.log(permissionsArray.current);
    dispatch(
      VerificationWindowActions.UpdatePermissionsAction({
        id: data.id,
        permissions: permissionsArray.current,
      }),
    );
  };
  return (
    <>
      {IsLoading === true ? (
        <GridLoading />
      ) : (
        <Wrapper className="permissionsWrapper">
          <FormControl component="fieldset">
            <FormLabel className="permissionsTitle" component="legend">
              Permissions
            </FormLabel>
            <div className="divider"></div>
            <FormGroup>
              {Permissions.map((item: Permission, index: number) => {
                return (
                  <div
                    key={'permission' + item.id}
                    className={`perm ${
                      index === Permissions.length - 1 ? 'last' : ''
                    }`}
                  >
                    <UbCheckbox
                      initialValue={item.userHasIt}
                      title={item.name}
                      titlePlacement={'end'}
                      onChange={checked => {
                        handleChange(index, checked, item.id);
                      }}
                    />
                  </div>
                );
              })}
            </FormGroup>
            <IsLoadingWithTextAuto
              icon={<SaveOutlinedIcon />}
              text="Save Changes"
              className={Buttons.SkyBlueButton}
              loadingId={'PermissionsButton' + data.id}
              onClick={handlePermissionSubmitClick}
            />
          </FormControl>
        </Wrapper>
      )}
    </>
  );
}

export default memo(PermissionsSegment);
const Wrapper = styled.div`
  width: 100%;
`;
