package graph

import "context"

// loadTagsForMeaning загружает теги для meaning.
func (r *Resolver) loadTagsForMeaning(ctx context.Context, meaningID int64) ([]*Tag, error) {
	meaningTags, err := r.tags.GetByMeaningIDs(ctx, []int64{meaningID})
	if err != nil {
		return nil, err
	}

	if len(meaningTags) == 0 {
		return []*Tag{}, nil
	}

	tagIDs := make([]int64, 0, len(meaningTags))
	for _, mt := range meaningTags {
		tagIDs = append(tagIDs, mt.TagID)
	}

	tags, err := r.tags.GetByIDs(ctx, tagIDs)
	if err != nil {
		return nil, err
	}

	return ToGraphQLTags(tags), nil
}

// ptrToInt безопасно разыменовывает указатель на int.
func ptrToInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}
