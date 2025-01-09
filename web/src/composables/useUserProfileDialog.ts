import { useDialog } from 'primevue/usedialog';
import { defineAsyncComponent } from 'vue';

/**
 *
 * @param user the user id, can be `me`
 */
export const useUserProfileDialog = function () {
  const dialog = useDialog();

  const UserProfileComponent = defineAsyncComponent(
    () => import('@/components/UserProfile.dialog.vue'),
  );

  return (user: string) =>
    dialog.open(UserProfileComponent, {
      props: {
        showHeader: false,
        closable: true,
        closeOnEscape: true,
        modal: true,
        style: {
          width: '50vw',
          'padding-top': '1rem',
          'overflow-y': 'scroll',
          'overflow-x': 'hidden',
        },
        breakpoints: {
          '960px': '75vw',
          '640px': '90vw',
        },
      },
      data: { user },
    });
};
