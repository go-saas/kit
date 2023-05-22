package price

import (
	"context"
	"github.com/go-saas/kit/pkg/localize"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"testing"
)

func TestPrice(t *testing.T) {
	p, err := NewPriceFromInt64(1000, "CNY")
	assert.NoError(t, err)
	ctx := localize.NewLanguageTagsContext(context.TODO(), []language.Tag{language.Chinese})
	pricePb := p.ToPricePb(ctx)
	assert.Equal(t, "¥10.00", pricePb.Text)
}
