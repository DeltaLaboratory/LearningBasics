import type { PageLoad } from './$types';

export const load: PageLoad = ({ params }) => {
    return {
        article_id: parseInt(params.id, 10)
    }
}